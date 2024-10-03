import React, { useEffect, useState } from "react";
import {
  PlayerUiState,
  PlayerUiContext,
  types,
  VideoClient,
  ControlBar,
  MediaContainer,
  PlayerAudioButton,
  PlayerBitrateButton,
  PlayerFullscreenButton,
  PlayerGetSoundButton,
  PlayerOverlayButton,
  PlayerPlayButton,
  PlayerVideo,
  PlayerVolumeRange,
} from "@video/video-client-web";
import { backendEndpoint } from "./utils.ts";

const ManifestPlayer = () => {
    const [vc, setVc] = useState<types.VideoClientAPI | null>(null);
    const [playerUi, setPlayerUi] = useState<PlayerUiState | null>(null);
    const [manifestUrl, setManifestUrl] = useState<string | null>(null);
    const [inputUrl, setInputUrl] = useState<string | "">("");
    useEffect(() => {
        if (vc == null) {
          const opts = {
            backendEndpoints: [backendEndpoint],
            userId: "demo",
          };
          const newVc = new VideoClient(opts);
          setVc(newVc);
        }
        return () => {
          if (vc != null) {
            vc.dispose();
            setVc(null);
          }
        };
      }, [vc]);

      useEffect(() => {
        if (vc != null && playerUi == null && manifestUrl) {
          const options = {};
          const player = vc.requestPlayer(manifestUrl, options);
          setPlayerUi(new PlayerUiState(player));
        }
        return () => {
          if (playerUi != null) {
            playerUi.dispose();
            setPlayerUi(null);
          }
        };
      }, [vc, playerUi, manifestUrl]);


      return (
        <>
        <h3>Please provide a valid manifest URL.</h3>
          <input
            type="text"
            placeholder="Enter manifest URL"
            value={inputUrl}
            onChange={(e) => setInputUrl(e.target.value)}
          />
          <button onClick={() => setManifestUrl(inputUrl)}>Set Manifest URL</button>
          {manifestUrl !== "" && playerUi && (
            <PlayerUiContext.Provider value={playerUi}>
              <MediaContainer>
                <PlayerGetSoundButton />
                <PlayerVideo />
                <ControlBar variant="player">
                  <PlayerPlayButton />
                  <PlayerAudioButton />
                  <PlayerVolumeRange />
                  <PlayerBitrateButton />
                  <PlayerFullscreenButton />
                </ControlBar>
                <PlayerOverlayButton />
              </MediaContainer>
            </PlayerUiContext.Provider>
          )}
        </>
      );
  };
  
  export default ManifestPlayer;