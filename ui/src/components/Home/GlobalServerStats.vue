<template>
  <v-card>
    <v-card-title class="dd-card-title dd-card-title-bottom-border">Global Stats</v-card-title>
    <v-simple-table dense>
      <tbody>
        <tr>
          <td class="dd-card-list-item-bold">Global Players</td>
          <td class="dd-card-list-item text-right">{{ data.global_players.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Game Time</td>
          <td class="dd-card-list-item text-right">{{ data.global_game_time.toLocaleString() }}s</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Deaths</td>
          <td class="dd-card-list-item text-right">{{ data.global_deaths.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Gems</td>
          <td class="dd-card-list-item text-right">{{ data.global_gems.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Enemies Killed</td>
          <td class="dd-card-list-item text-right">{{ data.global_enemies_killed.toLocaleString() }}</td>
        </tr>
        <tr @mouseover="globalAccuracyHover = true" @mouseleave="globalAccuracyHover = false">
          <td class="dd-card-list-item-bold">Global Accuracy</td>
          <td
            class="dd-card-list-item text-right"
            v-if="!globalAccuracyHover"
          >{{ data.global_accuracy }}%</td>
          <td class="dd-card-list-item text-right" v-else>
            {{ data.global_daggers_hit.toLocaleString() }} /
            {{ data.global_daggers_fired.toLocaleString() }}
          </td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Default Daggers</td>
          <td
            class="dd-card-list-item text-right"
          >{{ data.global_default_daggers.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Bronze Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.global_bronze_daggers.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Silver Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.global_silver_daggers.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Gold Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.global_gold_daggers.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Global Devil Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.global_devil_daggers.toLocaleString() }}</td>
        </tr>
      </tbody>
    </v-simple-table>
  </v-card>
</template>

<script>
export default {
  data() {
    return {
      globalAccuracyHover: false
    };
  },
  props: ["data"]
};
</script>

<style>
</style>