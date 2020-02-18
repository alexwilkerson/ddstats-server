<template>
  <v-card>
    <v-data-table
      :items="data.players"
      :loading="loading"
      :items-per-page="data.page_size"
      :page="data.page_number"
      :options.sync="options"
      :server-items-length="data.total_player_count"
      :hide-default-header="true"
      :disable-sort="true"
      :footer-props="{
      itemsPerPageOptions: [10],
      showFirstLastPage: true,
      showCurrentPage: true
    }"
      no-data-text="No players found."
      :mobile-breakpoint="NaN"
    >
      <template v-slot:header>
        <thead v-if="$root.mobile">
          <tr>
            <th class="text-left" title="Rank">
              <v-icon class="icon" color="#c33409" small>mdi-trophy</v-icon>
            </th>
            <th class="text-left" title="Player Name">
              <v-icon class="icon" color="#c33409" small>mdi-account</v-icon>
            </th>
            <th class="text-right" title="Game Time">
              <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>
            </th>
          </tr>
        </thead>
        <thead v-else>
          <tr>
            <th class="text-left pointer" @click="sort('rank')">Rank</th>
            <th class="text-left pointer" @click="sort('player_name')">Player Name</th>
            <th class="text-right pointer" @click="sort('game_time')">Highest Game Time</th>
            <th class="text-right pointer" @click="sort('overall_game_time')">Overall Game Time</th>
            <th class="text-right pointer" @click="sort('overall_deaths')">Overall Deaths</th>
            <th class="text-right pointer" @click="sort('overall_accuracy')">Overall Accuracy</th>
          </tr>
        </thead>
      </template>
      <template v-slot:body="{ items }">
        <tbody v-if="$root.mobile">
          <tr
            v-for="(item, i) in items"
            :key="i + item.player_id"
            @click="selectItem(item)"
            class="pointer"
          >
            <td class="grotesk rank">{{ item.rank }}</td>
            <td class="grotesk-bold">
              <v-icon :fill="$root.daggerColor(item.game_time)" small>$dagger</v-icon>
              {{ item.player_name }}
              <v-icon
                v-if="$root.checkPlayerLive(item.player_id)"
                class="icon online-green"
                small
              >mdi-access-point</v-icon>
            </td>
            <td
              class="text-right grotesk highest-game-time"
            >{{ Number.parseFloat(item.game_time).toFixed(4) }}</td>
          </tr>
        </tbody>
        <tbody v-else>
          <tr
            v-for="(item, i) in items"
            :key="i + item.player_id"
            @click="selectItem(item)"
            class="pointer"
          >
            <td class="grotesk rank">{{ item.rank }}</td>
            <td class="grotesk-bold">
              <v-icon :fill="$root.daggerColor(item.game_time)" small>$dagger</v-icon>
              {{ item.player_name }}
              <v-icon
                v-if="$root.checkPlayerLive(item.player_id)"
                class="icon online-green"
                small
              >mdi-access-point</v-icon>
            </td>
            <td
              class="text-right grotesk highest-game-time"
            >{{ Number.parseFloat(item.game_time).toFixed(4) }}s</td>
            <td
              class="text-right grotesk"
            >{{ moment.duration(item.overall_game_time, "seconds").humanize() }}</td>
            <td class="text-right grotesk" :style="{ width: '115px' }">{{ item.overall_deaths }}</td>
            <td
              class="text-right grotesk"
              :style="{ width: '135px' }"
            >{{ Number.parseFloat(item.overall_accuracy).toFixed(2) }}%</td>
          </tr>
        </tbody>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
const moment = require("moment");
import axios from "axios";
export default {
  data() {
    return {
      moment: moment,
      loading: true,
      sortBy: "rank",
      sortDir: "asc",
      data: {},
      options: {
        page: 1,
        rowsPerPage: 10
      }
    };
  },
  methods: {
    selectItem: function(item) {
      this.$router.push("/players/" + item.player_id);
    },
    getPlayersFromAPI() {
      this.loading = true;
      const { page, rowsPerPage } = this.options;
      axios
        .get(
          process.env.VUE_APP_API_URL +
            `/api/v2/player/all?page_size=${rowsPerPage}&page_num=${page}&sort_by=${this.sortBy}&sort_dir=${this.sortDir}`
        )
        .then(response => {
          this.data = response.data;
          this.loading = false;
        })
        .catch(error => window.console.log(error));
    },
    sort(by) {
      if (this.sortBy === by) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = by;
      }
      this.getPlayersFromAPI();
    }
  },
  mounted() {
    this.getPlayersFromAPI();
  },
  watch: {
    options: {
      handler() {
        this.getPlayersFromAPI();
      },
      deep: true
    }
  }
};
</script>

<style>
.v-data-table th {
  color: var(--v-primary-base) !important;
}
.red-text {
  color: var(--v-accent-base);
}
.grotesk {
  font-family: "alte_haas_grotesk", "Helvetica Neue", Helvetica, Arial;
}
.grotesk-bold {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
}
.pointer {
  cursor: pointer;
}
.rank {
  width: 35px;
}
.highest-game-time {
  width: 140px;
}
</style>
