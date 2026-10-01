import java.util.*;

// Ported from the owner's Greedy/JobSequence.java. Fixed: it took a job when
// its deadline was later than the number of jobs taken so far, so a rich job
// with a late deadline used up slot 1 and a job due at 1 was dropped ([2, 100]
// and [1, 50] gave 100, not 150). Each job now goes into the latest free slot
// at or before its deadline.
class Solution {
    static class Job {
        int deadline;
        int profit;
        int id;

        public Job(int i, int d, int p) {
            id = i;
            deadline = d;
            profit = p;
        }
    }

    public int maxProfit(int[][] jobsInfo) {
        ArrayList<Job> jobs = new ArrayList<>();
        for (int i = 0; i < jobsInfo.length; i++) {
            jobs.add(new Job(i, jobsInfo[i][0], jobsInfo[i][1]));
        }
        Collections.sort(jobs, (obj1, obj2) -> obj2.profit - obj1.profit);
        boolean[] used = new boolean[jobsInfo.length + 1];
        int total = 0;
        for (int i = 0; i < jobs.size(); i++) {
            Job curr = jobs.get(i);
            int slot = Math.min(curr.deadline, jobsInfo.length);
            while (slot > 0 && used[slot]) {
                slot--;
            }
            if (slot > 0) {
                used[slot] = true;
                total += curr.profit;
            }
        }
        return total;
    }
}
