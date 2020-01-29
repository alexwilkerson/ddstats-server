<template>
  <v-data-table :headers="headers" :options="options"></v-data-table>
</template>

<script>
export default {
  data: () => ({
    headers: [
      {
        text: "Rank",
        align: "left",
        sortable: false,
        value: "rank"
      },
      {
        text: "Player Name",
        sortable: false,
        value: "player_name"
      },
      {
        text: "Game Time",
        sortable: false,
        value: "game_time"
      },
      {
        text: "Gems",
        sortable: false,
        value: "gems"
      },
      {
        text: "Homing Daggers",
        sortable: false,
        value: "homing_daggers"
      },
      {
        text: "Accuracy",
        sortable: false,
        value: "accuracy"
      },
      {
        text: "Enemies Alive",
        sortable: false,
        value: "enemies_alive"
      },
      {
        text: "Enemies_killed",
        sortable: false,
        value: "enemies_killed"
      }
    ],
    options: {
      page: 1,
      itemsPerPage: 100
    }
  })
};
</script>

<style>
</style>