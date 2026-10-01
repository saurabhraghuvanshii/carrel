import java.util.*;

class Solution {
    public int[] mergeAll(int[][] lists) {
        // Each heap entry is {value, list, position}.
        PriorityQueue<int[]> heap = new PriorityQueue<>((a, b) -> Integer.compare(a[0], b[0]));
        int total = 0;
        for (int i = 0; i < lists.length; i++) {
            total += lists[i].length;
            if (lists[i].length > 0) {
                heap.add(new int[] {lists[i][0], i, 0});
            }
        }
        int[] out = new int[total];
        for (int w = 0; !heap.isEmpty(); w++) {
            int[] top = heap.poll();
            out[w] = top[0];
            int next = top[2] + 1;
            if (next < lists[top[1]].length) {
                heap.add(new int[] {lists[top[1]][next], top[1], next});
            }
        }
        return out;
    }
}
