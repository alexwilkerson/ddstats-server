<template>
  <div>
    <h2>Games Page. ID: {{ $route.params.id }}</h2>
    <GameCharts />
  </div>
</template>

<script>
import GameCharts from "../components/GameCharts";
export default {
  components: {
    GameCharts
  }
};
</script>
