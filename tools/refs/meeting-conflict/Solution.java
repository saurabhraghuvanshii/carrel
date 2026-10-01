import java.util.*;

class Solution {
    public boolean canAttendAll(int[][] meetings) {
        int[][] m = meetings.clone();
        Arrays.sort(m, (a, b) -> Integer.compare(a[0], b[0]));
        for (int i = 1; i < m.length; i++) {
            if (m[i][0] < m[i - 1][1]) {
                return false;
            }
        }
        return true;
    }
}
