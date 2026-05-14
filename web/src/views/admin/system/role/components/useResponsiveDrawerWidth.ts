import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

export function useResponsiveDrawerWidth(desktopWidth: number, tabletWidth = desktopWidth) {
  const windowWidth = ref(typeof window === 'undefined' ? 1440 : window.innerWidth)

  const syncWindowWidth = () => {
    windowWidth.value = window.innerWidth
  }

  onMounted(() => {
    syncWindowWidth()
    window.addEventListener('resize', syncWindowWidth)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('resize', syncWindowWidth)
  })

  const drawerWidth = computed<number | string>(() => {
    if (windowWidth.value < 768) {
      return '100vw'
    }
    if (windowWidth.value < 1280) {
      return Math.min(tabletWidth, windowWidth.value - 40)
    }
    return desktopWidth
  })

  return {
    drawerWidth
  }
}
