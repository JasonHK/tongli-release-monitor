import { z } from "zod";

// const NotifyBase = z.object({
//     with_title: z.boolean().default(true),
//     with_links: z.boolean().default(true),
// });

export type TelegramChatId = z.infer<typeof TelegramChatId>;
const TelegramChatId = z.union([z.string().startsWith("@"), z.number().int()])

const TelegramMessageThread = z.object({
    chat_id: TelegramChatId,
    message_thread_id: z.number().int(),
});

export type NotifyBase = z.infer<typeof NotifyBase>;
const NotifyBase = z.object({
    service: z.string(),
});

export type TelegramNotify = z.infer<typeof TelegramNotify>;
const TelegramNotify = NotifyBase.extend({
    service: z.literal("telegram"),
    bot_token: z.string(),
    chat_ids: z.union([TelegramChatId, TelegramMessageThread]).array().min(1),
});

export type Notify = z.infer<typeof Notify>;
export const Notify = z.discriminatedUnion("service", [
    TelegramNotify,
]);

export type Config = z.infer<typeof Config>;
export const Config = z.object({
    notify: z.array(Notify),
});
