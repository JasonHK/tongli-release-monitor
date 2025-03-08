import path from "node:path";
import fs from "node:fs";
import { createHash } from "node:crypto";
import { z } from "zod";
import { AsyncTask, SimpleIntervalJob, ToadScheduler } from "toad-scheduler";

import { Config } from "./config";
import { getLatestReleaseInfo, getReleaseBooks, ReleaseBooks } from "./scraper";

const LastScrapeInfo = z.object({
    url: z.string(),
    hash: z.string(),
});

class LastScrape
{
    readonly #path: string;

    #url: string = "";
    #hash: string = "";

    constructor(path: string)
    {
        this.#path = path;

        if (fs.existsSync(path))
        {
            const { url, hash } = LastScrapeInfo.parse(JSON.parse(fs.readFileSync(path, { encoding: "utf-8" })));
            this.#url = url;
            this.#hash = hash;
        }
    }

    /** Last scraped page's URL. */
    get url(): string { return this.#url; }
    
    /** Last scraped page's hash. */
    get hash(): string { return this.#hash; }

    /**
     * Update state when new page was scraped.
     * 
     * @param url  Scraped page's URL.
     * @param hash Scraped page's hash.
     */
    update(url: string, hash: string): void
    {
        this.#url = url;
        this.#hash = hash;

        fs.writeFileSync(this.#path, JSON.stringify({ url, hash }, null, 4));
    }
}

function calculateHash(books: ReleaseBooks): string
{
    const hash = createHash("md5");
    hash.update(JSON.stringify(Array.from(books).map(([k, v]) => ([k, Array.from(v)]))));

    return hash.digest("hex");
}

export function createMonitorTask(dataPath: string, config: Config)//: AsyncTask
{
    const lastScrape = new LastScrape(path.resolve(dataPath, "last_scrape"));

    const task = new AsyncTask("tongli-release-monitor", async () =>
    {
        const info = await getLatestReleaseInfo();
        const books = await getReleaseBooks(info.url);
        const hash = calculateHash(books);

        if (hash != lastScrape.hash)
        {
            lastScrape.update(info.url, hash);

            for (const target of config.notify)
            {
                const { default: notify } = await import(`./notifiers/${target.service}`);
                if (typeof notify === "function")
                {
                    notify({ ...info, books }, target);
                }
            }
        }
    });

    return task;
}
