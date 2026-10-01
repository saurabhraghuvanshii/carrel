import java.util.*;

class Solution {
    public int[] mergeOrder(int[][] logs, int[] offsets, int limit) {
        // {true time, server, position}
        PriorityQueue<int[]> heap = new PriorityQueue<>((a, b) -> a[0] != b[0] ? Integer.compare(a[0], b[0]) : Integer.compare(a[1], b[1]));
        for (int s = 0; s < logs.length; s++) {
            if (logs[s].length > 0) {
                heap.add(new int[] {logs[s][0] + offsets[s], s, 0});
            }
        }
        int[] out = new int[limit];
        for (int i = 0; i < limit; i++) {
            int[] e = heap.poll();
            int s = e[1], next = e[2] + 1;
            out[i] = s;
            if (next < logs[s].length) {
                heap.add(new int[] {logs[s][next] + offsets[s], s, next});
            }
        }
        return out;
    }
}
