// utils.ts
import { LoggerGlobal, LoggerCore } from "@video/log-client";
import { TokenRequest, types } from "@video/video-client-web";


export const backendEndpoint = "https://umbrella.dev1.devspace.lsea4.livelyvideo.tv";

/**
 * LoggerGlobal (only need one per application)
 * */
const loggerGlobal = new LoggerGlobal();
loggerGlobal.setOptions({
  host: "https://dev1.devspace.lsea4.livelyvideo.tv",
  interval: 5000,
  level: "debug",
});

/**
 * LoggerCore (only need one per application)
 * */
export const logger = new LoggerCore("VDC-web:BasicDemo");
logger.setLoggerMeta("client", "VDC");
logger.setLoggerMeta("chain", "VideoClient");
logger.setLoggerAggregate("message", "sample message");

/**
 * fetchToken
 * */
export const fetchToken = async (
  authUrl: string,
  reqBody: TokenRequest
): Promise<string> => {
  const response = await window.fetch(authUrl, {
    method: "post",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(reqBody),
    credentials: "include",
  });
  if (response.status !== 200) {
    throw new Error("Unable to get token");
  }

  const body = await response.json();
  return body.token;
};

export type TokenRefresherOptions = {
  authUrl: string;
  streamKey: string;
  scope: string;
  displayName?: string;
  userId?: string;
  clientReferrer?: string;
  streamName?: string;
};

/**
 * tokenRefresher
 * */
export const tokenRefresher =
  (options: TokenRefresherOptions): types.TokenGetter =>
  //This needs to be asynchronous because the fetchToken method will need to do a **POST** request to the authentication API
  async (): Promise<string> => {
    const url = `${options.authUrl}`;

    let token: string;
    try {
      const fetchOptions = {
        scopes: [options.scope],
        userId: options.userId ?? options.streamKey,
        data: {
          displayName: options.displayName ?? options.streamKey,
          mirrors: [
            {
              id: options.streamKey,
              streamName: options.streamName ?? "demo",
              kind: "pipe",
              clientEncoder: "demo",
              streamKey: options.streamKey,
              clientReferrer: options.clientReferrer ?? "demo",
            }
          ]
        },
      };
      token = await fetchToken(url, fetchOptions);
    } catch (error) {
      console.error("Error fetching token", error);
      throw error;
    }

    return token;
  };

  /* eslint-disable no-bitwise */
export function uuidv4(): string {
    return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
      const r = (Math.random() * 16) | 0;
      const v = c === "x" ? r : (r & 0x3) | 0x8;
      return v.toString(16);
    });
}
  