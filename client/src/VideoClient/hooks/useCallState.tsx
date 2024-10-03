// useCallState.tsx
import { CallState } from "@video/video-client-web";
import { useEffect, useState } from "react";

export function useCallState(): CallState | null {
  const [callState, setCallState] = useState<CallState | null>(null);

  /*
   * Create CallState.
   */
  useEffect(() => {
    if (callState == null) {
      setCallState(new CallState());
    }

    return () => {
      // Handle cleanup of CallState and Broadcast
      if (callState) {
        callState.stopBroadcast();
        callState.call?.close("Closed by call state on unmount/re-render");
        callState.dispose();
        setCallState(null);
      }
    };
  }, [callState]);

  return callState;
}