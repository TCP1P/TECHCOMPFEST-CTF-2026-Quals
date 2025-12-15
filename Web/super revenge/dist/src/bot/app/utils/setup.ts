/**
 * Copyright (c) 2025 replican
 * All rights reserved.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

import { chromium, firefox, webkit } from "playwright";
import type { Browser, BrowserContext, LaunchOptions } from "playwright";
import fs from "fs";
import path from "path";
import crypto from "crypto";

const browserHeadless = !(
  process.env.PLAYWRIGHT_DISPLAY === undefined &&
  fs.existsSync("/tmp/.X11-unix")
);


const USER_DATA_DIR = path.resolve(`/tmp/tmp-user-data-${Date.now()}-${crypto.randomBytes(4).toString("hex")}`);

let sharedBrowser: Browser | null = null;

/**
 * Returns a comma-separated string of full paths to extension directories.
 */


function getBrowserExtension(): string {
  const extDir = path.resolve("extensions"); // ensures full path
  if (!fs.existsSync(extDir)) return "";

  const dirs = fs.readdirSync(extDir).filter((file) => {
    const fullPath = path.join(extDir, file);
    return fs.lstatSync(fullPath).isDirectory();
  });

  const allExt = dirs
    .map((dir) => path.join(extDir, dir))
    .map((p) => path.resolve(p)) // ensure full absolute paths
    .join(",");

  console.info(`Loaded extensions: ${allExt}`);
  return allExt;
}

/**
 * Constructs launch options for the browser.
 * @param extension - Comma-separated list of extension directories.
 */
function getLaunchOptions(extension: string, browser: string): LaunchOptions {
  const baseArgs: string[] = [
    "--no-sandbox",
    "--disable-dev-shm-usage",
    "--disable-setuid-sandbox",
    "--disable-gpu",
    "--disable-background-networking",
    "--disable-background-timer-throttling",
    "--disable-renderer-backgrounding",
  ];

  if (browser === "chromium") {
    baseArgs.push("--disable-features=site-per-process");
  }

  const extensionArgs =
    extension && browser === "chromium"
      ? [
          `--disable-extensions-except=${extension}`,
          `--load-extension=${extension}`,
        ]
      : [
          "--disable-extensions",
          "--disable-translate",
          "--js-flags=--jitless,--noexpose_wasm",
        ];

  const options: LaunchOptions = {
    headless: browserHeadless,
    args: [...baseArgs, ...extensionArgs],
  };

  if (browser === "chromium") {
    options.channel = "chromium";
  }

  return options;
}

/**
 * Launches a browser based on the provided options.
 * Reuses a shared instance if no extensions are loaded.
 * @param options - Launch options for the browser.
 */
async function launchBrowser(
  options: LaunchOptions,
  browser: string
): Promise<Browser> {
  if (!getBrowserExtension() && sharedBrowser) {
    return sharedBrowser;
  }

  if (browser === "firefox") {
    return firefox.launch(options);
  } else {
    // Default to chromium
    return chromium.launch(options);
  }
}

async function launchBrowserWithExtension(
  options: LaunchOptions,
  browser: string
): Promise<BrowserContext> {
  if (browser === "chromium") {
    const context = chromium.launchPersistentContext(USER_DATA_DIR, options);
    return context;
  } else if (browser === "firefox") {
    const context = firefox.launchPersistentContext(USER_DATA_DIR, options);
    return context;
  } else if (browser === "webkit") {
    const context = webkit.launchPersistentContext(USER_DATA_DIR, options);
    return context;
  } else {
    throw new Error("Unsupported browser type");
  }
}

/**
 * Initializes and returns a new browser context with HTTPS errors ignored.
 * Uses a shared browser instance when no extensions are present.
 */
export async function browserCtx(
  browser: string,
  useExtension: boolean
): Promise<BrowserContext> {
  console.info(`Using browser type: ${browser}, headless: ${browserHeadless}`);
  if (browser === "firefox" && useExtension) {
    throw new Error("Firefox does not support extensions");
  }
  if (useExtension) {
    const extension = getBrowserExtension();
    const launchOptions = getLaunchOptions(extension, browser);
    if (extension) {
      return await launchBrowserWithExtension(launchOptions, browser);
    }
  }

  // No extension: use or reuse shared browser
  if (!sharedBrowser) {
    sharedBrowser = await launchBrowser(getLaunchOptions("", browser), browser);
  }

  return sharedBrowser.newContext({ ignoreHTTPSErrors: true });
}
