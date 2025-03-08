import process from "node:process";
import path from "node:path";
import fs from "node:fs";
import arg from "arg";
import { fromZodError, isZodErrorLike } from "zod-validation-error";
import { SimpleIntervalJob, ToadScheduler } from "toad-scheduler";
import YAML, { YAMLParseError } from "yaml";

import { Config } from "./config";
import { createMonitorTask } from "./monitor";

process.on("unhandledRejection", (reason) =>
{
    if (reason instanceof AggregateError)
    {
        reason.errors.forEach(console.error);
    }
    else
    {
        console.log(reason);
    }
});

const args = arg(
    {
        "--help": Boolean,
        "--config": String,
        "--data": String,
        "--once": Boolean,

        "-h": "--help",
        "-c": "--config",
        "-d": "--data",
    });

const configPath = args["--config"] ?? path.resolve(process.cwd(), "config.yml");
const dataPath = args["--data"] ?? path.resolve(process.cwd(), "data");

if (!fs.existsSync(configPath))
{
    console.error(`${configPath} does not exist.`);
    process.exit(1);
}
else if (!fs.statSync(configPath).isFile())
{
    console.error(`${configPath} must be a file.`);
    process.exit(2);
}

let config: Config;
try
{
    config = Config.parse(YAML.parse(fs.readFileSync(configPath, { encoding: "utf-8" })));
}
catch (error: unknown)
{
    if (error instanceof YAMLParseError)
    {
        console.error(error.message);
        process.exit(1);
    }
    else if (isZodErrorLike(error))
    {
        console.error(fromZodError(error).message);
        process.exit(1);
    }

    throw error;
}

if (!fs.existsSync(dataPath))
{
    fs.mkdirSync(dataPath);
}
else if (!fs.statSync(dataPath).isDirectory())
{
    console.error("--data must be a directory.");
    process.exit(1);
}

const task = createMonitorTask(dataPath, config);

if (args["--once"])
{
    task.execute();
}
else
{
    const scheduler = new ToadScheduler();
    const job = new SimpleIntervalJob({ hours: 1, runImmediately: true }, task, { preventOverrun: true });
    scheduler.addIntervalJob(job);

    process.on("SIGINT", () =>
    {
        process.stdout.moveCursor(-2, 0);
    
        console.info("Stopping running jobs...");
        scheduler.stop();
    });
}
