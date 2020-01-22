import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: Home
  },
  {
    path: "/games/:id",
    name: "games",
    component: () =>
      import(/* webpackChunkName: "games" */ "../views/Games.vue")
  },
  {
    path: "/games",
    name: "gamesList",
    component: () =>
      import(/* webpackChunkName: "gamesList" */ "../views/GamesList.vue")
  },
  {
    path: "/players",
    name: "playersList",
    component: () =>
      import(/* webpackChunkName: "playersList" */ "../views/PlayersList.vue")
  },
  {
    path: "/players/:id",
    name: "players",
    component: () =>
      import(/* webpackChunkName: "players" */ "../views/Players.vue")
  },
  {
    path: "/info",
    name: "info",
    component: () => import(/* webpackChunkName: "info" */ "../views/Info.vue")
  },
  {
    path: "/download",
    name: "download",
    component: () =>
      import(/* webpackChunkName: "download" */ "../views/Download.vue")
  },
  {
    path: "**",
    name: "notFound",
    component: () =>
      import(/* webpackChunkName: "notFound" */ "../views/NotFound.vue")
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
