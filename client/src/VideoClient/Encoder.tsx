// Encoder.tsx
import {
    CameraButton,
    ControlBar,
    EncoderVideo,
    JoinBroadcastButton,
    MediaContainer,
    MicrophoneButton,
    EncoderAudioDeviceSelect,
    EncoderResolutionSelect,
    EncoderVideoDeviceSelect,
    SettingsSidebar,
    TestMicButton,
    FullscreenButton,
    SettingsButton,
    VideoClientContext,
    EncoderUiContext,
    CallContext,
  } from "@video/video-client-web";
  import React from "react";
import { useVideoClient } from "./hooks/useVideoClient.tsx";
import { useEncoderUi } from "./hooks/useEncoderUi.tsx";
import { useCallState } from "./hooks/useCallState.tsx";
  
  function Encoder(): React.ReactElement {
    /**
     * Access `VideoClient`, `EncoderUiState`, and
     * `CallState` from your custom hooks.
     * */
    const videoClient = useVideoClient();
    const encoderUi = useEncoderUi();
    const callState = useCallState();
  
    /** NOTE: Do not interact with EncoderUiContext
     * or VideoClientContext in this component, as it
     * would be OUTSIDE of the EncoderUiProvider and
     * VideoClientProvider. This component is only for rendering.
     * */
  
    return (
      <VideoClientContext.Provider value={videoClient}>
        <CallContext.Provider value={callState}>
          <EncoderUiContext.Provider value={encoderUi}>
            {encoderUi != null && callState != null && (
              <>
                {/* MediaContainer should wrap 
                ALL components for styling. */}
                <MediaContainer>
                  <EncoderVideo />
                  {/* ControlBar wraps controls (for styling). 
                  Include required variant prop */}
                  <ControlBar variant="encoder">
                    <JoinBroadcastButton
                      broadcastOptions={{
                        streamName: "demo",
                      }}
                    />
                    <CameraButton />
                    <MicrophoneButton />
                    <FullscreenButton />
                    <SettingsButton />
                  </ControlBar>
                  {/* SettingsSidebar wraps items to be 
                  displayed in sidebar (for styling). */}
                  <SettingsSidebar>
                    <div>
                      <EncoderVideoDeviceSelect />
                      <EncoderAudioDeviceSelect />
                      <EncoderResolutionSelect />
                    </div>
                    <TestMicButton />
                  </SettingsSidebar>
                </MediaContainer>
              </>
            )}
          </EncoderUiContext.Provider>
        </CallContext.Provider>
      </VideoClientContext.Provider>
    );
  }
  export default Encoder;