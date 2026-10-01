import java.util.*;

class Solution {
    public int earliestArrival(int n, int[] period, int[][] roads) {
        List<List<int[]>> out = new ArrayList<>();
        for (int i = 0; i < n; i++) {
            out.add(new ArrayList<>());
        }
        for (int[] r : roads) {
            out.get(r[0]).add(r);
        }
        int[] best = new int[n];
        Arrays.fill(best, Integer.MAX_VALUE);
        best[0] = 0;
        PriorityQueue<int[]> pq = new PriorityQueue<>((a, b) -> Integer.compare(a[1], b[1]));
        pq.add(new int[] {0, 0});
        while (!pq.isEmpty()) {
            int[] cur = pq.poll();
            int u = cur[0];
            if (cur[1] > best[u]) {
                continue;
            }
            int leave = (cur[1] + period[u] - 1) / period[u] * period[u];
            for (int[] r : out.get(u)) {
                if (leave + r[2] < best[r[1]]) {
                    best[r[1]] = leave + r[2];
                    pq.add(new int[] {r[1], best[r[1]]});
                }
            }
        }
        return best[n - 1] == Integer.MAX_VALUE ? -1 : best[n - 1];
    }
}
