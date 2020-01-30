<template>
  <v-data-table
    :items="data.games"
    :headers="this.$root.mobile ? mobileHeaders : headers"
    :loading="loadingTable"
    :items-per-page="data.page_size"
    :page="data.page_num"
    :server-items-length="data.total_game_count"
    :disable-sort="true"
    :footer-props="{
      itemsPerPageOptions: [10],
      showFirstLastPage: true,
      showCurrentPage: true
    }"
    :mobile-breakpoint="NaN"
  >
    <template v-slot:body="{ items }">
      <tbody v-if="$root.mobile">
        <tr v-for="item in items" :key="item.player_id" @click="selectItem(item)" class="pointer">
          <td class="text-right grotesk">{{ item.rank }}</td>
          <td class="grotesk-bold red-text">{{ item.player_name }}</td>
          <td class="text-right grotesk">{{ Number.parseFloat(item.game_time).toFixed(4) }}</td>
        </tr>
      </tbody>
      <tbody v-else>
        <tr v-for="item in items" :key="item.player_id" @click="selectItem(item)" class="pointer">
          <td class="text-right grotesk">{{ item.rank }}</td>
          <td class="grotesk-bold red-text">{{ item.player_name }}</td>
          <td class="text-right grotesk">{{ Number.parseFloat(item.game_time).toFixed(4) }}</td>
          <td class="text-right grotesk">{{ item.gems }}</td>
          <td class="text-right grotesk">{{ item.homing_daggers }}</td>
          <td class="text-right grotesk">{{ Number.parseFloat(item.accuracy).toFixed(2) }}%</td>
          <td class="text-right grotesk">{{ item.enemies_alive }}</td>
          <td class="text-right grotesk">{{ item.enemies_killed }}</td>
        </tr>
      </tbody>
    </template>
  </v-data-table>
</template>

<script>
export default {
  props: ["data", "loadingTable"],
  data: () => ({
    headers: [
      {
        text: "Rank",
        align: "right",
        value: "rank"
      },
      {
        text: "Player Name",
        align: "left",
        value: "player_name"
      },
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
      }
    ],
    mobileHeaders: [
      {
        text: "Rank",
        align: "right",
        value: "rank"
      },
      {
        text: "Player Name",
        align: "left",
        value: "player_name"
      },
      {
        text: "Game Time",
        align: "right",
        value: "game_time"
      }
    ]
  }),
  methods: {
    selectItem: function(item) {
      this.$router.push("/games/" + item.id);
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
</style>
