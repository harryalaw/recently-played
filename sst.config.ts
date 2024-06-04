import { SSTConfig } from "sst";
import { Cron, Function } from "sst/constructs";


export default {
    config(input) {
        return {
            name: "spotify-tracker",
            region: "eu-west-2",
            profile: input.stage === "production" ? "prod" : "dev"
        };
    },
    stacks(app) {
        app.setDefaultFunctionProps({
            runtime: "go",
        });
        app.stack(function Stack({ stack }) {
            const recently_played = new Function(stack, "recently-played-sst", {
                handler: "./cmd/functions/recently-played/",
                environment: {
                    SPOTIFY_CLIENT_ID: process.env.SPOTIFY_CLIENT_ID,
                    SPOTIFY_CLIENT_SECRET: process.env.SPOTIFY_CLIENT_SECRET,
                    SPOTIFY_REFRESH_TOKEN: process.env.SPOTIFY_REFRESH_TOKEN,
                    DSN: process.env.DSN
                }
            });

            const song_of_the_day = new Function(stack, "song-of-the-day-sst", {
                handler: "./cmd/functions/song-of-the-day/",
                environment: {
                    SPOTIFY_CLIENT_ID: process.env.SPOTIFY_CLIENT_ID,
                    SPOTIFY_CLIENT_SECRET: process.env.SPOTIFY_CLIENT_SECRET,
                    SPOTIFY_REFRESH_TOKEN: process.env.SPOTIFY_REFRESH_TOKEN,
                    PLAYLIST_ID: process.env.PLAYLIST_ID,
                    DSN: process.env.DSN
                }
            });

            new Cron(stack, "recently-played-trigger", {
                job: recently_played,
                schedule: "rate(30 minutes)",
                enabled: !app.local
            });

            new Cron(stack, "song-of-the-day-trigger", {
                job: song_of_the_day,
                schedule: "rate(1 day)",
                enabled: !app.local
            });
        });


    }
} satisfies SSTConfig
