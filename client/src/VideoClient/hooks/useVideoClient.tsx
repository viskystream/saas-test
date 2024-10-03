// useVideoClient.tsx
import { types, VideoClient } from "@video/video-client-web";
import { useEffect, useState } from "react";
import { logger, uuidv4, tokenRefresher, backendEndpoint } from "../utils.ts"

export function useVideoClient(): types.VideoClientAPI | null {
  const [videoClient, setVideoClient] = useState<types.VideoClientAPI | null>(
    null
  );
  useEffect(() => {
    if (videoClient == null) {
      const token = tokenRefresher({
        
        authUrl: `${backendEndpoint}/apps/demos/api/demo/v1/access-token`,
        scope: "conference-owner",
        streamKey: uuidv4(),
      });

      /**
       * Setting the generated token and the backendEndpoint for
       * the options to be passed to our new VideoClient instance
       **/
      const videoClientOptions: types.VideoClientOptions = {
        backendEndpoints: [backendEndpoint],
        token,
        userId: "demo",
        logger,
      };

      const vc = new VideoClient(videoClientOptions);
      setVideoClient(vc);
    }

    return () => {
      // Handle cleanup of VideoClient instance
      if (videoClient != null) {
        videoClient.dispose();
        setVideoClient(null);
      }
    };
    /*
     * Remember to only include things in your dependency array
     * that are related to the state of your `VideoClient` instance,
     * otherwise disposal may occur at undesired times.
     */
  }, [videoClient]);

  return videoClient;
}