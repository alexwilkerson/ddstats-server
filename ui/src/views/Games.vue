<template>
  <div class="wrapper">
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

<style scoped>
.wrapper {
  padding-top: 40px;
  padding-bottom: 40px;
}
</style>
