import java.util.*;

class Solution {
    public int maxRiders(int stops, int seats, int[][] rides) {
        int[][] byDropOff = rides.clone();
        Arrays.sort(byDropOff, (a, b) -> Integer.compare(a[1], b[1]));
        int[] load = new int[stops];
        int taken = 0;
        for (int[] r : byDropOff) {
            boolean fits = true;
            for (int x = r[0]; x < r[1] && fits; x++) {
                fits = load[x] < seats;
            }
            if (fits) {
                for (int x = r[0]; x < r[1]; x++) {
                    load[x]++;
                }
                taken++;
            }
        }
        return taken;
    }
}
