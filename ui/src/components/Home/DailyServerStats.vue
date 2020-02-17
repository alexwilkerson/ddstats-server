<template>
  <v-card>
    <v-card-title class="dd-card-title dd-card-title-bottom-border">Today's Stats</v-card-title>
    <v-simple-table dense>
      <tbody>
        <tr>
          <td class="dd-card-list-item-bold">Game Time</td>
          <td class="dd-card-list-item text-right">{{ data.since_game_time.toLocaleString() }}s</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Deaths</td>
          <td class="dd-card-list-item text-right">{{ data.since_deaths.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Gems</td>
          <td class="dd-card-list-item text-right">{{ data.since_gems.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Enemies Killed</td>
          <td class="dd-card-list-item text-right">{{ data.since_enemies_killed.toLocaleString() }}</td>
        </tr>
        <tr @mouseover="sinceAccuracyHover = true" @mouseleave="sinceAccuracyHover = false">
          <td class="dd-card-list-item-bold">Accuracy</td>
          <td
            class="dd-card-list-item text-right"
            v-if="!sinceAccuracyHover"
          >{{ data.since_accuracy }}%</td>
          <td class="dd-card-list-item text-right" v-else>
            {{ data.since_daggers_hit.toLocaleString() }} /
            {{ data.since_daggers_fired.toLocaleString() }}
          </td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">New Players</td>
          <td class="dd-card-list-item text-right">{{ data.new_players.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Active Players</td>
          <td class="dd-card-list-item text-right">{{ data.active_players.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">New Devil Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.since_devil_daggers.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">New Gold Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.since_gold_daggers.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">New Silver Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.since_silver_daggers.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">New Bronze Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.since_bronze_daggers.toLocaleString() }}</td>
        </tr>
      </tbody>
    </v-simple-table>
  </v-card>
</template>

<script>
export default {
  data() {
    return {
      sinceAccuracyHover: false
    };
  },
  props: ["data"]
};
</script>

<style>
</style>