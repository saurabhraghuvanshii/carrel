import java.util.*;

class Solution {
    public int roomsNeeded(int[][] meetings) {
        int[][] m = meetings.clone();
        Arrays.sort(m, (a, b) -> Integer.compare(a[0], b[0]));
        PriorityQueue<Integer> ends = new PriorityQueue<>();
        for (int[] meeting : m) {
            if (!ends.isEmpty() && ends.peek() <= meeting[0]) {
                ends.poll();
            }
            ends.add(meeting[1]);
        }
        return ends.size();
    }
}
