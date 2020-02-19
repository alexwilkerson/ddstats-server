<template>
  <v-card>
    <v-data-table
      :items="data.games"
      :loading="loadingTable"
      :items-per-page="data.page_size"
      :server-items-length="data.total_game_count"
      :headers="headers"
      :options.sync="options"
      :disable-sort="true"
      :hide-default-header="true"
      :footer-props="{
        itemsPerPageOptions: [10],
        showFirstLastPage: true,
        showCurrentPage: true,
        disablePagination: loadingTable
      }"
      no-data-text="No leaderboard found."
      :mobile-breakpoint="NaN"
    >
      <template v-slot:header>
        <thead v-if="$root.mobile">
          <tr>
            <th></th>
            <th
              class="text-left pointer"
              title="Player Name"
              @click="sort('player_name')"
            >
              <v-icon class="icon" color="#c33409" small>mdi-account</v-icon>
            </th>
            <th
              class="text-right pointer"
              title="Game Time"
              @click="sort('game_time')"
            >
              <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>
            </th>
          </tr>
        </thead>
        <thead v-if="!$root.mobile">
          <tr>
            <th class="text-left pointer" @click="sort('rank')">Rank</th>
            <th class="text-left pointer" @click="sort('player_name')">
              Player Name
            </th>
            <th class="text-right pointer" @click="sort('game_time')">
              Game Time
            </th>
            <th class="text-right pointer" @click="sort('gems')">Gems</th>
            <th class="text-right pointer" @click="sort('homing_daggers')">
              Homing Daggers
            </th>
            <th class="text-right pointer" @click="sort('accuracy')">
              Accuracy
            </th>
            <th class="text-right pointer" @click="sort('enemies_alive')">
              Enemies Alive
            </th>
            <th class="text-right pointer" @click="sort('enemies_killed')">
              Enemies Killed
            </th>
          </tr>
        </thead>
      </template>
      <template v-slot:body="{ items }">
        <tbody v-if="$root.mobile">
          <tr
            v-for="item in items"
            :key="item.player_id"
            @click="selectItem(item)"
            class="pointer"
          >
            <td class="text-right grotesk">{{ item.rank }}</td>
            <td class="grotesk-bold">
              <v-icon :fill="$root.daggerColor(item.game_time)" small
                >$dagger</v-icon
              >
              {{ item.player_name }}
              <v-icon
                v-if="$root.checkPlayerLive(item.player_id)"
                class="icon online-green"
                small
                >mdi-access-point</v-icon
              >
            </td>
            <td class="text-right grotesk">
              {{ Number.parseFloat(item.game_time).toFixed(4) }}
            </td>
          </tr>
        </tbody>
        <tbody v-else>
          <tr
            v-for="item in items"
            :key="item.player_id"
            @click="selectItem(item)"
            class="pointer"
          >
            <td class="text-right grotesk">{{ item.rank }}</td>
            <td class="grotesk-bold">
              <v-icon
                v-if="item.game_time >= data.devil_dagger_time"
                fill="#c33409"
                small
                >$dagger</v-icon
              >
              <v-icon
                v-else-if="item.game_time >= data.gold_dagger_time"
                fill="#ffcd00"
                small
                >$dagger</v-icon
              >
              <v-icon
                v-else-if="item.game_time >= data.silver_dagger_time"
                fill="#acacac"
                small
                >$dagger</v-icon
              >
              <v-icon
                v-else-if="item.game_time >= data.bronze_dagger_time"
                fill="#ff8300"
                small
                >$dagger</v-icon
              >
              <v-icon v-else fill="#000" small>$dagger</v-icon>
              {{ item.player_name }}
              <v-icon
                v-if="$root.checkPlayerLive(item.player_id)"
                class="icon online-green"
                small
                >mdi-access-point</v-icon
              >
            </td>
            <td class="text-right grotesk">
              {{ Number.parseFloat(item.game_time).toFixed(4) }}
            </td>
            <td class="text-right grotesk">{{ item.gems }}</td>
            <td class="text-right grotesk">{{ item.homing_daggers }}</td>
            <td class="text-right grotesk">
              {{ Number.parseFloat(item.accuracy).toFixed(2) }}%
            </td>
            <td class="text-right grotesk">{{ item.enemies_alive }}</td>
            <td class="text-right grotesk">{{ item.enemies_killed }}</td>
          </tr>
        </tbody>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
export default {
  props: ["data", "loadingTable", "sort", "optionsChanged"],
  data: () => ({
    options: {
      page: 1,
      rowsPerPage: 10
    },
    headers: [{}, {}, {}, {}, {}, {}, {}, {}] // this is required to make the stupid loading progress bar the correct length
  }),
  methods: {
    selectItem: function(item) {
      this.$router.push("/games/" + item.id);
    }
  },
  watch: {
    options: {
      handler() {
        this.$emit("optionsChanged", this.options);
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
.leaderboard-dagger {
  padding: 0px !important;
}
</style>
