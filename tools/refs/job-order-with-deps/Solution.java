import java.util.*;

class Solution {
    public int finishTime(int n, int[] durations, int[][] deps) {
        List<List<Integer>> after = new ArrayList<>();
        for (int i = 0; i < n; i++) {
            after.add(new ArrayList<>());
        }
        int[] waiting = new int[n];
        for (int[] d : deps) {
            after.get(d[0]).add(d[1]);
            waiting[d[1]]++;
        }
        int[] start = new int[n];
        ArrayDeque<Integer> ready = new ArrayDeque<>();
        for (int j = 0; j < n; j++) {
            if (waiting[j] == 0) {
                ready.add(j);
            }
        }
        int done = 0, finish = 0;
        while (!ready.isEmpty()) {
            int j = ready.poll();
            done++;
            int end = start[j] + durations[j];
            finish = Math.max(finish, end);
            for (int next : after.get(j)) {
                start[next] = Math.max(start[next], end);
                if (--waiting[next] == 0) {
                    ready.add(next);
                }
            }
        }
        return done == n ? finish : -1;
    }
}
