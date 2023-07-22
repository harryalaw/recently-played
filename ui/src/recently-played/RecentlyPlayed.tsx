/*
 *
 * What do I want to do here?
 * I want to get all of the data from the DB, then in the UI I want to get the relevant spotify data
 *
 * The boring way is to just get the tracks in batches
 *  -> I think to rethink this we would need to store artist/album data instead to more easily batch on that
 *
 *  -> Just show a list of track IDs to start with
 *
 *
 */

import { useQuery } from "@tanstack/react-query";
import { TsvData } from "./parse-local-data";
import { usePlayedCounts } from "./use-played-counts";

export function RecentlyPlayed() {
  const { isLoading, isError, data, error } = usePlayedCounts("2");
  if (isLoading) {
    return <div> Loading ... </div>;
  }
  if (isError) {
    console.error(error);
    return <div> Error! </div>;
  }

  return <PlayedList hrefs={data.map((el) => el.href)} />;
}

function PlayedList(props: { hrefs: string[] }) {
  const { isLoading, isError, data, error } = useSpotifyTrackData(props.hrefs);

  if (isLoading) {
    return <div> Loading ... </div>;
  }
  if (isError) {
    console.error(error);
    return <div> Error! </div>;
  }

  return (
    <ol>
      {data.map((track, id) => {
        console.log(track, id);
        if (track !== undefined) {
          return (
            <li key={id}>
              {track.name}
              {track.artists.map((artist) => artist.name).join(" ")}
              {track.album.name}
            </li>
          );
        }
      })}
    </ol>
  );
}

function useSpotifyTrackData(hrefs: string[]) {
  return useQuery({
    queryKey: ["tracks"],
    queryFn: () => getSpotifyTrackData(hrefs),
    staleTime: Infinity,
  });
}

async function getSpotifyTrackData(hrefs: string[]) {
  const uniqueHrefs = getUnique(hrefs);
  const trackIds = uniqueHrefs.map((href) => href.split("/")[5]);

  const uniqueCount = trackIds.length;
  console.log(uniqueCount);

  const trackInfo: Record<string, TrackInfo> = {};

  // batch into 50s
  // limiting to 50 to reduce spotify requests while testing!
  const batchSize = 50;
  console.log(batchSize);
  for (let i = 0; i < uniqueCount; i += batchSize) {
    const section = trackIds.slice(i, i + batchSize);
    console.log(section);
    const tracks = await requestSpotifyTrackData(section);
    for (const track of tracks.tracks) {
      trackInfo[track.id ?? "UNKNOWN"] = track;
    }
  }

  console.log(trackInfo);

  const data = hrefs
    .map((href) => href.split("/")[5])
    .map((trackId) => trackInfo[trackId]);

  return data;
}

type TrackInfo = {
  album: AlbumInfo;
  artists: ArtistInfo[];
  available_markets?: string[];
  disc_number?: number;
  duration_ms?: number;
  explicit?: boolean;
  external_ids?: ExternalIds;
  external_urls?: ExternalUrls;
  href?: string;
  id?: string;
  is_playable?: boolean;
  linked_from?: any;
  restrictions?: Restrictions;
  name?: string;
  popularity?: number;
  preview_url?: string;
  track_number?: number;
  type?: "track";
  uri?: string;
  is_local?: boolean;
};

type Restrictions = {
  reason?: "market" | "product" | "explicit";
};

type ExternalIds = {
  isrc?: string;
  ean?: string;
  upc?: string;
};

type ExternalUrls = {
  spotify?: string;
};

type AlbumInfo = {
  album_type: string;
  total_tracks: number;
  external_urls: ExternalUrls;
  href: string;
  id: string;
  images: ImageObject[];
  name: string;
  release_date: string;
  release_data_precision: "year" | "month" | "day";
  restrictions?: Restrictions;
  type: "album";
  uri: string;
  copyrights?: {
    text?: string;
    type?: string;
  };
  external_ids?: ExternalIds;
  genres?: string[];
  label?: string;
  popularity?: number;
  album_group?: "album" | "single" | "compilation" | "appears_on";
  artists: SimpleArtistInfo[];
};

type ImageObject = {
  url: string;
  height: number | null;
  width: number | null;
};

type SimpleArtistInfo = Pick<
  ArtistInfo,
  "external_urls" | "href" | "id" | "name" | "type" | "uri"
>;

type ArtistInfo = {
  external_urls?: {
    spotify?: string;
  };
  followers?: {
    href: null;
    total?: number;
  };
  genres?: string[];
  href?: string;
  id?: string;
  images?: ImageObject[];
  name?: string;
  popularity?: number;
  type?: "artist";
  uri?: string;
};

async function requestSpotifyTrackData(trackIds: string[]) {
  const authToken = localStorage.getItem("access_token");
  let apiUrl = `https://api.spotify.com/v1/tracks?ids=${trackIds.join(",")}`;

  console.log(apiUrl);
  console.log(apiUrl.length);

  const headers = new Headers();
  headers.append("Authorization", `Bearer ${authToken}`);

  return fetch(apiUrl, {
    method: "GET",
    headers: { Authorization: `Bearer ${authToken}` },
  }).then((response) => {
    if (!response.ok) {
      throw new Error("Network response was not ok: " + response.status);
    }
    return response.json() as Promise<{ tracks: TrackInfo[] }>;
  });
}

function getUnique<T>(items: T[]): T[] {
  return [...new Set(items)];
}
