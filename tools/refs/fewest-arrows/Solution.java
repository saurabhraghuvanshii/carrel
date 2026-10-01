import java.util.*;

class Solution {
    public int fewestArrows(int[][] balloons) {
        int[][] b = balloons.clone();
        Arrays.sort(b, (x, y) -> Integer.compare(x[1], y[1]));
        int arrows = 0;
        long at = Long.MIN_VALUE;
        for (int[] balloon : b) {
            if (balloon[0] > at) {
                arrows++;
                at = balloon[1];
            }
        }
        return arrows;
    }
}
