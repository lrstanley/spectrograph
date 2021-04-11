<template>
    <div>
        <div class="header text-center">
            <h1>It's time to use <span class="bravery--text">Spectrograph</span> in your server.</h1>
            <p class="text-subtitle-1">
                Spectrograph is a Discord bot that manages creating/destroying
                <br> voice channels depending on occupancy, so you don't have to.
            </p>
            <v-btn color="primary" @click.up="anchor('more-info')">More info</v-btn>
            <v-btn color="success">Add to server</v-btn>
        </div>
        <span class="chevron-down d-flex justify-center">
            <v-icon size="75">mdi-chevron-down</v-icon>
        </span>
        <v-container id="more-info" ref="more-info" class="px-7 px-md-0">
            <!-- <div style="height: 3000px">Hello world test</div> -->
            <v-row justify="center">
                <v-col cols="12" md="2" class="d-flex justify-center">
                    <v-img contain max-height="400" src="/src/static/img/example-bad.png" />
                </v-col>
                <v-col cols="12" md="6" class="d-flex flex-column justify-lg-center">
                    <h1>Out with the old.</h1>
                    <p class="text-subtitle-1">
                        Often times, larger (and some smaller) servers have these types of
                        channels setup. This can get annoying -- so many channels that are empty.
                        And what about smaller servers that don't have this kind of setup, that
                        run out of space? Why not be more efficient, without requiring mods/admins
                        to step in and create spare channels?

                        <br>
                        <br>
                        <strong>Spectrograph does just that</strong>. The idea is to always have
                        one empty channel for a given voice channel "group". If there are 2, or
                        500 voice channels, there is always an available one to join, ensuring
                        there is always an adequate amount. No less, no more.
                    </p>
                </v-col>
            </v-row>
            <v-row justify="center">
                <v-col cols="12" md="6" class="d-flex flex-column justify-lg-center">
                    <h1>In with the new (and simple).</h1>
                    <p class="text-subtitle-1">
                        Spectrograph monitors voice channel member activity, as well as the
                        creation/deletion/editing of channels.
                        <ul>
                            <li>
                                Any time there are more than 1 empty channels for a given voice
                                channel group, Spectrograph will remove the extra channels,
                                ensuring the empty channel is always at the bottom.
                            </li>
                            <li>
                                Any time there are no empty channels (e.g. user just joined the
                                last empty channel), Spectrograph will create a new channel at
                                the bottom of the voice channel group, copying the configuration
                                from the other channels in the group.
                            </li>
                            <li>
                                When the first voice channel in the given group is changed
                                (bitrate, user limit, permissions, etc), Spectrograph will
                                persist these changes across all other channels of the voice
                                channel group.
                            </li>
                        </ul>
                    </p>
                </v-col>
                <v-col cols="12" md="2" class="d-flex justify-center">
                    <v-img contain max-height="400" src="/src/static/img/example-new.png" />
                </v-col>
            </v-row>
        </v-container>
    </div>
</template>

<script>
export default {
    name: "index",
    title: "Index",
    data: function () {
        return {}
    },
}
</script>

<style scoped>
.header {
    /* padding-top: 200px; */
    padding: 200px 30px 0;
    height: 130vh;
    margin-bottom: -300px;
    width: 100vw;
    color: #ffffff;
    background-image: url(/src/static/img/bg.svg);
    background-position: center;
    background-size: 200%;
    background-repeat: no-repeat;
}

/* @media only screen and (max-width: 600px) {}
@media only screen and (min-width: 600px) {}
@media only screen and (min-width: 768px) {}
@media only screen and (min-width: 992px) {} */

@media only screen and (min-width: 1200px) {
    .header {
        /* padding-top: 200px; */
        background-position: center; /* top; */
        background-size: 100%;
    }
}

.chevron-down {
    position: relative;
    top: -15vh;
}

#more-info {
    padding-top: 100px;
    margin-bottom: 100px;
}

#more-info > .row {
    margin-bottom: 50px;
}
</style>
