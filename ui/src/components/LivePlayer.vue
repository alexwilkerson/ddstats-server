<template>
  <div>
    game time: {{ $root.state.game_time.toFixed(4) }}s <br />
    gems: {{ $root.state.gems }} <br />
    homing daggers: {{ $root.state.homing_daggers }} <br />
    enemies alive: {{ $root.state.enemies_alive }} <br />
    enemies killed: {{ $root.state.enemies_killed }} <br />
    accuracy:
    <span v-if="$root.state.daggers_fired > 0">
      {{
        (($root.state.daggers_hit / $root.state.daggers_fired) * 100).toFixed(
          2
        )
      }}%
    </span>
    <span v-else>0%</span>
    <br />
    level two time: {{ $root.state.level_two_time }} <br />
    level three time: {{ $root.state.level_three_time }} <br />
    level four time: {{ $root.state.level_four_time }} <br />
    death type:
    {{ getDeathType($root.state.death_type) }}
    <br />
    is replay: {{ $root.state.is_replay }} <br />

    status:
    {{
      $root.state.death_type !== -1 && $root.status === "Alive"
        ? "Dead"
        : $root.status
    }}
  </div>
</template>

<script>
export default {
  data() {
    return {
      state: {}
    };
  },
  methods: {
    getDeathType(id) {
      if (id === -1) {
        return "N/A";
      }
      return deathTypes[id];
    }
  }
};
const deathTypes = [
  "FALLEN",
  "SWARMED",
  "IMPALED",
  "GORED",
  "INFESTED",
  "OPENED",
  "PURGED",
  "DESECRATED",
  "SACRIFICED",
  "EVISCERATED",
  "ANNIHILATED",
  "INTOXICATED",
  "ENVENMONATED",
  "INCARNATED",
  "DISCARNATED",
  "BARBED"
];
</script>

<style></style>
