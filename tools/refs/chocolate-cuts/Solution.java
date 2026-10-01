import java.util.*;

// Ported from the owner's Greedy/ChocolaPrblm.java (written as a main). Fixed:
// the total was an int, which overflows on large bars; it is now a long.
class Solution {
    public long cutCost(int n, int m, int[] hor, int[] ver) {
        Integer costVer[] = new Integer[ver.length];
        Integer costHor[] = new Integer[hor.length];
        for (int i = 0; i < ver.length; i++) {
            costVer[i] = ver[i];
        }
        for (int i = 0; i < hor.length; i++) {
            costHor[i] = hor[i];
        }
        Arrays.sort(costVer, Collections.reverseOrder());
        Arrays.sort(costHor, Collections.reverseOrder());
        int h = 0, v = 0;
        int hp = 1, vp = 1;
        long cost = 0;
        while (h < costHor.length && v < costVer.length) {
            if (costVer[v] <= costHor[h]) {
                cost += (long) costHor[h] * vp;
                hp++;
                h++;
            } else {
                cost += (long) costVer[v] * hp;
                vp++;
                v++;
            }
        }
        while (h < costHor.length) {
            cost += (long) costHor[h] * vp;
            hp++;
            h++;
        }
        while (v < costVer.length) {
            cost += (long) costVer[v] * hp;
            vp++;
            v++;
        }
        return cost;
    }
}
