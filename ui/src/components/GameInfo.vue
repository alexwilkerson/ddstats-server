<template>
  <div class="wrapper">
    <v-card :style="{padding: '12px'}">
      <v-row no-gutters justify="center">
        <v-col cols="12" sm="2" align="center">
          <v-tooltip bottom nudgeRight="6" nudgeTop="6" contentClass="tooltip">
            <template v-slot:activator="{ on }">
              <div v-on="on" class="icon-info">
                <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>
                <span>{{ gameInfo.game_time }}s</span>
              </div>
            </template>
            <span>Game Time</span>
          </v-tooltip>
        </v-col>
        <v-col cols="12" sm="2" align="center">
          <v-tooltip bottom nudgeRight="5" nudgeTop="6" contentClass="tooltip">
            <template v-slot:activator="{ on }">
              <div v-on="on" class="icon-info">
                <v-icon class="icon" fill="#c33409" small>$gem</v-icon>
                <span>{{ gameInfo.gems }}</span>
              </div>
            </template>
            <span>Gems</span>
          </v-tooltip>
        </v-col>
        <v-col cols="12" sm="2" align="center">
          <v-tooltip bottom nudgeRight="4" nudgeTop="6" contentClass="tooltip">
            <template v-slot:activator="{ on }">
              <div v-on="on" class="icon-info">
                <v-icon class="icon" fill="#c33409">$dagger</v-icon>
                <span>{{ gameInfo.homing_daggers }}</span>
              </div>
            </template>
            <span>Homing Daggers</span>
          </v-tooltip>
        </v-col>
        <v-col cols="12" sm="2" align="center">
          <v-tooltip bottom nudgeRight="2" nudgeTop="6" contentClass="tooltip">
            <template v-slot:activator="{ on }">
              <div v-on="on" class="icon-info">
                <v-icon class="icon" fill="#c33409" small>$target</v-icon>
                <span>{{ gameInfo.accuracy }}%</span>
              </div>
            </template>
            <span>Accuracy</span>
          </v-tooltip>
        </v-col>
        <v-col cols="12" sm="2" align="center">
          <v-tooltip bottom nudgeRight="4" nudgeTop="6" contentClass="tooltip">
            <template v-slot:activator="{ on }">
              <div v-on="on" class="icon-info">
                <v-icon class="icon" fill="#c33409" small>$skull</v-icon>
                <span>{{ gameInfo.enemies_alive }}</span>
              </div>
            </template>
            <span>Enemies Alive</span>
          </v-tooltip>
        </v-col>
        <v-col cols="12" sm="2" align="center">
          <v-tooltip bottom nudgeRight="4" nudgeTop="6" contentClass="tooltip">
            <template v-slot:activator="{ on }">
              <div v-on="on" class="icon-info">
                <v-icon class="icon" fill="#c33409" small>$splat</v-icon>
                <span>{{ gameInfo.enemies_killed }}</span>
              </div>
            </template>
            <span>Enemies Killed</span>
          </v-tooltip>
        </v-col>
      </v-row>
      <hr :style="{margin: '12px 0'}" />
      <v-simple-table dense>
        <template v-slot:default>
          <tbody>
            <tr
              @click="$router.push('/players/' + gameInfo.player_id)"
              :style="{ cursor: 'pointer' }"
            >
              <td class="text-left">
                <v-icon
                  class="icon"
                  :style="{ marginLeft: '4px', marginRight: '10px' }"
                  color="#c33409"
                  small
                >mdi-account</v-icon>Player
              </td>
              <td class="text-right">{{ gameInfo.player_name }}</td>
            </tr>
            <tr>
              <td class="text-left">
                <v-icon
                  class="icon"
                  :style="{ marginLeft: '4px', marginRight: '10px' }"
                  color="#c33409"
                  small
                >mdi-card-bulleted</v-icon>Game ID
              </td>
              <td class="text-right">{{ $route.params.id }}</td>
            </tr>
            <tr>
              <td class="text-left">
                <v-icon
                  class="icon"
                  :style="{ marginLeft: '4px', marginRight: '10px' }"
                  color="#c33409"
                  small
                >mdi-grave-stone</v-icon>Death Type
              </td>
              <td class="text-right">{{ gameInfo.death_type }}</td>
            </tr>
            <tr>
              <td class="text-left">
                <v-icon
                  class="icon"
                  :style="{ marginLeft: '4px', marginRight: '10px' }"
                  color="#c33409"
                  small
                >mdi-sim</v-icon>Spawnset
              </td>
              <td class="text-right">{{ gameInfo.spawnset }}</td>
            </tr>
            <tr @mouseover="dateHover = true" @mouseleave="dateHover = false">
              <td class="text-left">
                <v-icon
                  class="icon"
                  :style="{ marginLeft: '4px', marginRight: '10px' }"
                  color="#c33409"
                  small
                >mdi-calendar-month</v-icon>Recorded
              </td>
              <td class="text-right">
                {{
                dateHover
                ? moment
                .utc(gameInfo.time_stamp)
                .local()
                .format("lll")
                : moment(gameInfo.time_stamp).fromNow()
                }}
              </td>
            </tr>
            <tr
              @click="$router.push('/players/' + gameInfo.replay_player_id)"
              v-if="gameInfo.replay_player_name !== undefined"
              :style="{ cursor: 'pointer' }"
            >
              <td class="text-left">
                <v-icon
                  class="icon"
                  :style="{ marginLeft: '4px', marginRight: '10px' }"
                  color="#c33409"
                  small
                >mdi-video-account</v-icon>Replay Recorded By
              </td>
              <td class="text-right">{{ gameInfo.replay_player_name }}</td>
            </tr>
          </tbody>
        </template>
      </v-simple-table>
    </v-card>
  </div>
</template>

<script>
const moment = require("moment");
export default {
  data() {
    return {
      moment: moment,
      nameHover: false,
      dateHover: false
    };
  },
  props: {
    gameInfo: {
      type: Object,
      required: true
    }
  }
};
</script>

<style scoped>
tr:hover {
  /* background: #fffefc !important; */
  /* background: var(--v-highlight-base) !important; */
}
.stats-table {
  border-radius: 2px;
  max-width: 650px;
  margin: 0 auto 10px auto;
}
.stats-table tr {
  /* background: #f6f2ee; */
  /* background: var(--v-background-base); */
}
.wrapper {
  max-width: 860px;
  margin: 0 auto 20px auto;
  /* background: #f6f2ee; */
  /* background: var(--v-background-base); */
}
.icon {
  margin-top: -4px;
  margin-right: 6px;
}
.game-info {
  padding: 0;
}
.game-info li {
  font-size: 15px;
  display: inline;
  cursor: default;
}
.game-info div {
  display: inline;
}
.icon-info span {
  font-size: 15px;
  cursor: default;
}
.tooltip {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 12px;
  border-radius: 2px;
}
tbody {
  /* background: #ebe7e4; */
}
ul {
  list-style: none;
  font-size: 14px;
}
</style>
