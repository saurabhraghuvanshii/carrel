import java.util.*;

class Inventory {
    private final Map<Integer, Integer> free = new HashMap<>();
    private final Map<Integer, int[]> held = new HashMap<>(); // order -> {item, qty}

    public void add(int item, int qty) {
        free.merge(item, qty, Integer::sum);
    }

    public boolean reserve(int order, int item, int qty) {
        if (held.containsKey(order) || free.getOrDefault(item, 0) < qty) {
            return false;
        }
        free.merge(item, -qty, Integer::sum);
        held.put(order, new int[] {item, qty});
        return true;
    }

    public boolean cancel(int order) {
        int[] h = held.remove(order);
        if (h == null) {
            return false;
        }
        free.merge(h[0], h[1], Integer::sum);
        return true;
    }

    public boolean confirm(int order) {
        return held.remove(order) != null;
    }

    public int available(int item) {
        return free.getOrDefault(item, 0);
    }
}
