import java.util.*;

class Solution {
    public int[][] insert(int[][] intervals, int[] add) {
        List<int[]> out = new ArrayList<>();
        int i = 0;
        int n = intervals.length;
        while (i < n && intervals[i][1] < add[0]) {
            out.add(intervals[i++]);
        }
        int start = add[0];
        int end = add[1];
        while (i < n && intervals[i][0] <= end) {
            start = Math.min(start, intervals[i][0]);
            end = Math.max(end, intervals[i][1]);
            i++;
        }
        out.add(new int[] {start, end});
        while (i < n) {
            out.add(intervals[i++]);
        }
        return out.toArray(new int[0][]);
    }
}
