<template>
  <v-card>
    <v-data-table
      :items="data.games"
      :headers="this.$root.mobile ? mobileHeaders : headers"
      :loading="loading"
      :items-per-page="data.page_size"
      :page="data.page_number"
      :options.sync="options"
      :server-items-length="data.total_game_count"
      :disable-sort="true"
      :hide-default-header="true"
      :footer-props="{
        itemsPerPageOptions: [10],
        showFirstLastPage: true,
        showCurrentPage: true,
        disablePagination: loading
      }"
      no-data-text="No games found."
      :mobile-breakpoint="NaN"
    >
      <template v-slot:header>
        <thead v-if="$root.mobile">
          <tr>
            <th
              class="text-right pointer"
              title="Game Time"
              @click="sort('game_time')"
            >
              <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>
            </th>
            <th class="text-right pointer" title="Recorded" @click="sort('id')">
              <v-icon class="icon" color="#c33409" small
                >mdi-calendar-month</v-icon
              >
            </th>
          </tr>
        </thead>
        <thead v-else>
          <tr>
            <th
              class="text-right pointer"
              title="Game Time"
              @click="sort('game_time')"
            >
              Game Time
            </th>
            <th class="text-right pointer" title="Gems" @click="sort('gems')">
              Gems
            </th>
            <th
              class="text-right pointer"
              title="Homing Daggers"
              @click="sort('homing_daggers')"
            >
              Homing Daggers
            </th>
            <th
              class="text-right pointer"
              title="Accuracy"
              @click="sort('accuracy')"
            >
              Accuracy
            </th>
            <th
              class="text-right pointer"
              title="Enemies Alive"
              @click="sort('enemies_alive')"
            >
              Enemies Alive
            </th>
            <th
              class="text-right pointer"
              title="Enemies Killed"
              @click="sort('enemies_killed')"
            >
              Enemies Killed
            </th>
            <th class="text-right pointer" title="Recorded" @click="sort('id')">
              Recorded
            </th>
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
            <td class="text-right grotesk game-time">
              {{ Number.parseFloat(item.game_time).toFixed(4) }}
            </td>
            <td class="text-right grotesk recorded">
              {{ moment(item.time_stamp).fromNow() }}
            </td>
          </tr>
        </tbody>
        <tbody v-else>
          <tr
            v-for="(item, i) in items"
            :key="i + item.player_id"
            @click="selectItem(item)"
            class="pointer"
          >
            <td class="text-right grotesk game-time">
              {{ Number.parseFloat(item.game_time).toFixed(4) }}
            </td>
            <td class="text-right grotesk">{{ item.gems }}</td>
            <td class="text-right grotesk">{{ item.homing_daggers }}</td>
            <td class="text-right grotesk">
              {{ Number.parseFloat(item.accuracy).toFixed(2) }}%
            </td>
            <td class="text-right grotesk">{{ item.enemies_alive }}</td>
            <td class="text-right grotesk">{{ item.enemies_killed }}</td>
            <td class="text-right grotesk recorded">
              {{ moment(item.time_stamp).fromNow() }}
            </td>
          </tr>
        </tbody>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
const moment = require("moment");
import axios from "axios";
import EventBus from "../event-bus";

export default {
  data() {
    return {
      moment: moment,
      loading: true,
      data: {},
      sortBy: "id",
      sortDir: "desc",
      options: {
        page: 1,
        rowsPerPage: 10
      },
      headers: [
        {
          text: "Game Time",
          align: "right",
          value: "game_time"
        },
        {
          text: "Gems",
          align: "right",
          value: "gems"
        },
        {
          text: "Homing Daggers",
          align: "right",
          value: "homing_daggers"
        },
        {
          text: "Accuracy",
          align: "right",
          value: "accuracy"
        },
        {
          text: "Enemies Alive",
          align: "right",
          value: "enemies_alive"
        },
        {
          text: "Enemies Killed",
          align: "right",
          value: "enemies_killed"
        },
        {
          text: "Recorded",
          align: "right",
          value: "time_stamp"
        }
      ],
      mobileHeaders: [
        {
          text: "Game Time",
          align: "right",
          value: "game_time"
        },
        {
          text: "Recorded",
          align: "right",
          value: "time_stamp"
        }
      ]
    };
  },
  methods: {
    selectItem: function(item) {
      this.$router.push("/games/" + item.id);
    },
    getGamesFromAPI() {
      this.loading = true;
      const { page, rowsPerPage } = this.options;
      axios
        .get(
          process.env.VUE_APP_API_URL +
            `/api/v2/game/recent?player_id=${this.$route.params.id}&page_size=${rowsPerPage}&page_num=${page}&sort_by=${this.sortBy}&sort_dir=${this.sortDir}`
        )
        .then(response => {
          this.data = response.data;
          this.$emit("onPlayerNameLoad", response.data.player_name);
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
      this.getGamesFromAPI();
    }
  },
  mounted() {
    // this.getGamesFromAPI();
    EventBus.$on(
      "game_submitted",
      function(body) {
        if (body.player_id == this.$route.params.id) {
          this.getGamesFromAPI();
        }
      }.bind(this)
    );
  },
  watch: {
    options: {
      handler() {
        this.getGamesFromAPI();
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
.recorded {
  width: 120px;
}
.game-time {
  width: 80px;
}
</style>
