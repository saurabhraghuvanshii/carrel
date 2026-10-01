import java.util.*;

class Solution {
    public int[] sessions(int[][] events, int timeout) {
        Map<Integer, List<Integer>> byUser = new HashMap<>();
        for (int[] e : events) {
            byUser.computeIfAbsent(e[0], u -> new ArrayList<>()).add(e[1]);
        }
        int count = 0, longest = 0;
        for (List<Integer> times : byUser.values()) {
            Collections.sort(times);
            int start = times.get(0);
            count++;
            for (int i = 1; i < times.size(); i++) {
                if (times.get(i) - times.get(i - 1) > timeout) {
                    count++;
                    start = times.get(i);
                }
                longest = Math.max(longest, times.get(i) - start);
            }
        }
        return new int[] {count, longest};
    }
}
