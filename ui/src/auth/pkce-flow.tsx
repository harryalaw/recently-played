import { Route } from "@tanstack/router";
import { rootRoute } from "../main";
import { SpotifyStartAuth, SpotifyAuth } from "./AuthComponents";
import { z } from "zod";
import { RecentlyPlayedDashboard } from "../recently-played/RecentlyPlayed";

const VALID_CHARS =
  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

export function generateRandomString(length: number) {
  let text = "";

  for (let i = 0; i < length; i++) {
    text += VALID_CHARS[Math.floor(Math.random() * VALID_CHARS.length)];
  }

  return text;
}

function base64encode(string: ArrayBuffer): string {
  return btoa(
    String.fromCharCode.apply(
      null,
      new Uint8Array(string) as unknown as number[],
    ),
  )
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/g, "");
}

export async function generateCodeChallenge(
  codeVerifier: string,
): Promise<string> {
  const encoder = new TextEncoder();

  const data = encoder.encode(codeVerifier);
  const digest = await window.crypto.subtle.digest("SHA-256", data);

  return base64encode(digest);
}

export const initiateAuthRoute = new Route({
  getParentRoute: () => rootRoute,
  path: "/initiate",
  component: SpotifyStartAuth,
});

const spotifyAuthSchema = z.object({
  code: z.string().default(""),
  state: z.string().default(""),
});

export const receiveAuthRoute = new Route({
  getParentRoute: () => rootRoute,
  path: "/callback",
  component: SpotifyAuth,
  validateSearch: spotifyAuthSchema,
});

export const dashboardAuthRoute = new Route({
  getParentRoute: () => rootRoute,
  path: "/dashboard",
  component: RecentlyPlayedDashboard,
});
