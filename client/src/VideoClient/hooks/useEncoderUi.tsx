// useEncoderUi.tsx
import { EncoderUiState, mediaController } from "@video/video-client-web";
import { useEffect, useState } from "react";

export function useEncoderUi(): EncoderUiState | null {
  const [encoderUi, setEncoderUi] = useState<EncoderUiState | null>(null);

  /*
   * Create MediaStreamController + EncoderUiState for Broadcaster.
   */
  useEffect(() => {
    if (encoderUi == null) {
      (async () => {
        await mediaController.init();

        const mediaStreamController = await mediaController.requestController();
        setEncoderUi(new EncoderUiState(mediaStreamController));
      })();
    }
    return () => {
      // Handle cleanup of mediaStreamController and EncoderUiState
      if (encoderUi != null) {
        encoderUi.mediaStreamController?.close(
          "Closed by unmounting/re-render"
        );
        encoderUi.dispose("Component unmounting/re-render");
        setEncoderUi(null);
      }
    };
  }, [encoderUi]);

  return encoderUi;
}