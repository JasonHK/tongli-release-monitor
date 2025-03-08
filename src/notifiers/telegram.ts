import TelegramBot from "node-telegram-bot-api";

import { TelegramNotify } from "src/config";
import { Notifier, Release } from "src/notifier";

const ESCAPED_CHARS = /([_*[\]()~`>#+-=|{}.!])/g;

function escape(text: string): string
{
    return text.replaceAll(ESCAPED_CHARS, "\\$1");
}

async function notify({ title, url, books }: Release, { bot_token, chat_ids }: TelegramNotify)
{
    const bot = new TelegramBot(bot_token);

    let message = `[*${escape(title)}*](${url})`;
    for (const [category, bookNames] of books)
    {
        message += `\n\n__*${escape(category)}*__`;
        for (const book of bookNames)
        {
            message += `\n◦${escape(book)}`;
        }
    }

    for (const id of chat_ids)
    {
        let chat_id: TelegramBot.ChatId;
        let message_thread_id: number | undefined;
        if (typeof id === "object")
        {
            chat_id = id.chat_id;
            message_thread_id = id.message_thread_id;
        }
        else
        {
            chat_id = id;
        }

        bot.sendMessage(chat_id, message, { message_thread_id, parse_mode: "MarkdownV2", disable_web_page_preview: true });
    }
}

export default notify satisfies Notifier<TelegramNotify>;
