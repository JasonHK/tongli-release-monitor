import { NotifyBase } from "./config";
import { NewsInfo, ReleaseBooks } from "./scraper";

export interface Release extends NewsInfo
{
    books: ReleaseBooks;
}

export type Notifier<T extends NotifyBase> = (release: Release, config: T) => void | Promise<void>;
