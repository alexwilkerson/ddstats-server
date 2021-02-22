<template>
  <div class="wrapper">
    <h1>Info</h1>
    <div v-if="!loading">
      <Info />
    </div>
  </div>
</template>

<script>
import Info from "../components/Info";

export default {
  data: () => ({}),
  components: {
    Info
  }
};
</script>

<style scoped>
.wrapper {
  padding-left: 10px;
  padding-right: 10px;
  padding-top: 40px;
  padding-bottom: 40px;
  max-width: 860px;
  margin: auto;
}
h1 {
  text-align: center;
  margin-bottom: 15px;
  color: var(--v-primary-base);
}
</style>