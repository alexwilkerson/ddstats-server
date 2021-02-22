<template>
  <v-card>
    <v-toolbar flat color="th" dark dense>
      <v-toolbar-title>Today's Stats</v-toolbar-title>
    </v-toolbar>
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
          <td class="dd-card-list-item-bold">Players with New Scores</td>
          <td
            class="dd-card-list-item text-right"
          >{{ data.players_with_new_scores.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Players With New Ranks</td>
          <td
            class="dd-card-list-item text-right"
          >{{ data.players_with_new_ranks.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Avg. Improvement Time</td>
          <td class="dd-card-list-item text-right">{{ data.average_improvement_time.toFixed(4) }}s</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Avg. Rank Improvement</td>
          <td
            class="dd-card-list-item text-right"
          >{{ data.average_rank_improvement.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Avg. Game Time</td>
          <td
            class="dd-card-list-item text-right"
          >{{ data.average_game_time_per_active_player.toFixed(4) }}s</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Avg. Deaths</td>
          <td
            class="dd-card-list-item text-right"
          >{{ data.average_deaths_per_active_player.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Avg. Gems</td>
          <td
            class="dd-card-list-item text-right"
          >{{ data.average_gems_per_active_player.toLocaleString() }}</td>
        </tr>
        <tr>
          <td class="dd-card-list-item-bold">Avg. Enemies Killed</td>
          <td
            class="dd-card-list-item text-right"
          >{{ data.average_enemies_killed_per_active_player.toLocaleString() }}</td>
        </tr>
        <tr @mouseover="avgAccuracyHover = true" @mouseleave="avgAccuracyHover = false">
          <td class="dd-card-list-item-bold">Avg. Accuracy</td>
          <td
            class="dd-card-list-item text-right"
            v-if="!avgAccuracyHover"
          >{{ data.average_accuracy_per_active_player }}%</td>
          <td class="dd-card-list-item text-right" v-else>
            {{ data.average_daggers_hit_per_active_player.toLocaleString() }} /
            {{ data.average_daggers_fired_per_active_player.toLocaleString() }}
          </td>
        </tr>
        <tr v-if="$root.mobile">
          <td class="dd-card-list-item-bold">New Players</td>
          <td class="dd-card-list-item text-right">{{ data.new_players.toLocaleString() }}</td>
        </tr>
        <tr v-if="$root.mobile">
          <td class="dd-card-list-item-bold">Active Players</td>
          <td class="dd-card-list-item text-right">{{ data.active_players.toLocaleString() }}</td>
        </tr>
        <tr v-if="$root.mobile">
          <td class="dd-card-list-item-bold">New Devil Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.since_devil_daggers.toLocaleString() }}</td>
        </tr>
        <tr v-if="$root.mobile">
          <td class="dd-card-list-item-bold">New Gold Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.since_gold_daggers.toLocaleString() }}</td>
        </tr>
        <tr v-if="$root.mobile">
          <td class="dd-card-list-item-bold">New Silver Daggers</td>
          <td class="dd-card-list-item text-right">{{ data.since_silver_daggers.toLocaleString() }}</td>
        </tr>
        <tr v-if="$root.mobile">
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
      sinceAccuracyHover: false,
      avgAccuracyHover: false
    };
  },
  props: ["data"]
};
</script>

<style>
</style>