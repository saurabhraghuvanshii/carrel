import java.util.*;

class Solution {
    public int[][] freeTime(int[][][] busy) {
        // Walk all lists together with a heap of {start, end, person, position}.
        PriorityQueue<int[]> heap = new PriorityQueue<>((a, b) -> Integer.compare(a[0], b[0]));
        for (int p = 0; p < busy.length; p++) {
            heap.add(new int[] {busy[p][0][0], busy[p][0][1], p, 0});
        }
        List<int[]> free = new ArrayList<>();
        int end = heap.peek()[0];
        while (!heap.isEmpty()) {
            int[] top = heap.poll();
            if (top[0] > end) {
                free.add(new int[] {end, top[0]});
            }
            end = Math.max(end, top[1]);
            int next = top[3] + 1;
            if (next < busy[top[2]].length) {
                int[] iv = busy[top[2]][next];
                heap.add(new int[] {iv[0], iv[1], top[2], next});
            }
        }
        return free.toArray(new int[0][]);
    }
}
