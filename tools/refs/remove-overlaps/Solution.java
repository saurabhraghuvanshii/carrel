import java.util.*;

// Ported from the owner's Greedy/GreedyA.java (activity selection): it counts
// the most intervals that can be kept, sorted by end with start >= last end;
// the answer is everything else.
class Solution {
    public int fewestRemovals(int[][] intervals) {
        int[][] activities = new int[intervals.length][3];
        for (int i = 0; i < intervals.length; i++) {
            activities[i][0] = i;
            activities[i][1] = intervals[i][0];
            activities[i][2] = intervals[i][1];
        }
        Arrays.sort(activities, Comparator.comparingDouble(o -> o[2]));
        int maxAct = 1;
        int lastEnd = activities[0][2];
        for (int i = 1; i < intervals.length; i++) {
            if (activities[i][1] >= lastEnd) {
                maxAct++;
                lastEnd = activities[i][2];
            }
        }
        return intervals.length - maxAct;
    }
}
