import { JSDOM } from "jsdom";
import innerText from "styleless-innertext";

export interface NewsInfo
{
    title: string;
    url: string;
}

export type ReleaseBooks = Map<string, Set<string>>;

export async function getLatestReleaseInfo(): Promise<NewsInfo>
{
    const { window } = await JSDOM.fromURL("https://www.tongli.com.tw/TNews_List.aspx?Type=0&Page=1");

    const news = window.document.querySelector<HTMLLIElement>(".news_list li:first-child");
    if (!news) { throw new Error("No news was found. Possibly some changes in the UI."); }

    const title = news.querySelector<HTMLHeadingElement>(".title");
    const url = news.querySelector<HTMLAnchorElement>(".title a")?.href;
    if (!title || !url) { throw new Error("No news info was found. Possibly some changes in the UI."); }

    return { title: innerText(title), url };
}

export async function getReleaseBooks(url: string): Promise<ReleaseBooks>
{
    const { window } = await JSDOM.fromURL(url);

    const article = window.document.querySelector("#ContentPlaceHolder1_TNewsContent > div");
    if (!article) { throw new Error("No article was found. Possibly some changes in the UI."); }

    const books: ReleaseBooks = new Map();

    let category: string | undefined;
    for (const element of article.children)
    {
        if (element instanceof window.HTMLDivElement)
        {
            category = innerText(element);
            if (!books.has(category)) { books.set(category, new Set()); }
        }
        else if (element instanceof window.HTMLParagraphElement)
        {
            if (!category) { throw new Error("No category was set. Possibly some changes in the UI."); }
            
            for (const node of element.childNodes)
            {
                if (node instanceof window.Text)
                {
                    const title = node.wholeText.trim();
                    if (title) { books.get(category)!.add(node.wholeText.trim()); }
                }
            }
        }
    }

    return books;
}
