import { useQuery } from "@tanstack/react-query";
import { usePlayedCounts } from "./use-played-counts";
import { TsvData } from "./parse-local-data";

export function RecentlyPlayedDashboard() {
  const { isLoading, isError, data, error } = usePlayedCounts("2");
  if (isLoading) {
    return <div> Loading ... </div>;
  }
  if (isError) {
    console.error(error);
    return <div> Error! </div>;
  }

  return <PlayedList trackData={data} />;
}

function groupByDate(trackData: TrackData[]): Record<string, TrackData[]> {
  const dateMap: Record<string, TrackData[]> = {};

  for (const track of trackData) {
    const date = track.date.toDateString();
    const tracks = dateMap[date] ?? [];
    tracks.push(track);
    dateMap[date] = tracks;
  }

  for (const [_date, tracks] of Object.entries(dateMap)) {
    tracks.sort((a, b) => b.date.getTime() - a.date.getTime());
  }

  return dateMap;
}

function PlayedList(props: { trackData: TsvData[] }) {
  const { isLoading, isError, data, error } = useSpotifyTrackData(
    props.trackData,
  );

  if (isLoading) {
    return <div> Loading ... </div>;
  }
  if (isError) {
    console.error(error);
    return <div> Error! </div>;
  }

  const groupedByDate = groupByDate(data);
  const days = Object.entries(groupedByDate).sort((a, b) => {
    const dateA = new Date(a[0]);
    const dateB = new Date(b[0]);
    return dateB.getTime() - dateA.getTime();
  });

  return (
    <div>
      {days.map(([day, tracks]) => {
        return (
          <div>
            <div
              style={{
                position: "sticky",
                top: "0",
                background: "black",
                color: "white",
                fontSize: '1.5rem',
                padding: '1rem'
              }}
            >
              {day}
            </div>
            {tracks.map((track) => {
              return (
                <TrackDetails trackData={track} key={track.date.getTime()} />
              );
            })}
          </div>
        );
      })}
    </div>
  );
}

function msToMinuteSeconds(ms: number): string {
  const inSeconds = Math.floor(ms / 1000);
  const minutes = Math.floor(inSeconds / 60);
  const seconds = inSeconds % 60;

  return `${minutes}:${padTime(seconds)}`;
}

function padTime(time: number): string {
  return String(time).padStart(2, "0");
}

function TrackDetails(props: { trackData: TrackData }) {
  const { track, date } = props.trackData;

  return (
    <div
      role="row"
      style={{
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        gap: "1em",
        maxWidth: 800,
      }}
    >
      <div>{`${padTime(date.getHours())}:${padTime(date.getMinutes())}`}</div>
      <div>
        {
          <img
            src={track.album.images[2].url}
            alt={`Album artwork for ${track.album.name}`}
          />
        }
      </div>
      <div>
        <p>{track.name}</p>
        <p>{track.artists.map((artist) => artist.name).join(", ")}</p>
      </div>
      <div style={{ marginLeft: "auto" }}>
        <p>{msToMinuteSeconds(track.duration_ms ?? 0)}</p>
      </div>
    </div>
  );
}

function useSpotifyTrackData(trackData: TsvData[]) {
  return useQuery({
    queryKey: ["tracks"],
    queryFn: () => getSpotifyTrackData(trackData),
    staleTime: Infinity,
  });
}

type TrackData = {
  track: TrackInfo;
  date: Date;
};

async function getSpotifyTrackData(trackData: TsvData[]): Promise<TrackData[]> {
  const hrefs = trackData.map((track) => track.href);

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

  const data = trackData.map((data) => {
    const trackId = data.href.split("/")[5];
    const spotifyTrackData = trackInfo[trackId];
    return {
      track: spotifyTrackData,
      date: data.date,
    };
  });

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
