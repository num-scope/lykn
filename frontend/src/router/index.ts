import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "../stores/auth";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: () => import("../layouts/DefaultLayout.vue"),
      children: [
        {
          path: "",
          redirect: { name: "projects" },
        },
        {
          path: "login",
          name: "login",
          component: () => import("../views/LoginView.vue"),
        },
        {
          path: "projects",
          name: "projects",
          component: () => import("../views/ProjectsView.vue"),
        },
        {
          path: "licenses",
          name: "licenses",
          component: () => import("../views/LicensesView.vue"),
        },
        {
          path: "features",
          name: "features",
          component: () => import("../views/FeaturesView.vue"),
        },
        {
          path: "plans",
          name: "plans",
          component: () => import("../views/PlansView.vue"),
        },
      ],
    },
  ],
});

router.beforeEach((to) => {
  const authStore = useAuthStore();
  if (!authStore.isAuthenticated && to.name !== "login") {
    return { name: "login" };
  }
  if (authStore.isAuthenticated && to.name === "login") {
    return { name: "projects" };
  }
  return true;
});

export default router;
