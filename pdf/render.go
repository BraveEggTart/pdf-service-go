package pdf

import (
	"context"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// RenderHTMLToPDF 将 HTML 字符串渲染为 PDF 字节。ctx 应带超时（如 30s）。
func RenderHTMLToPDF(ctx context.Context, html string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)
	if path := os.Getenv("CHROME_PATH"); path != "" {
		opts = append(opts, chromedp.ExecPath(path))
	} else if path := os.Getenv("CHROMEDP_CHROME_PATH"); path != "" {
		opts = append(opts, chromedp.ExecPath(path))
	} else {
		// Docker 常见路径
		for _, p := range []string{"/usr/bin/chromium", "/usr/bin/chromium-browser", "/headless-shell/headless-shell"} {
			if _, err := os.Stat(p); err == nil {
				opts = append(opts, chromedp.ExecPath(p))
				break
			}
		}
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	var buf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	}

	if err := chromedp.Run(browserCtx, tasks); err != nil {
		return nil, err
	}
	return buf, nil
}
