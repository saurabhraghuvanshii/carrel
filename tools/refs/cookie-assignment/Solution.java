import java.util.*;

class Solution {
    public int feedChildren(int[] greed, int[] cookies) {
        Arrays.sort(greed);
        Arrays.sort(cookies);
        int fed = 0;
        for (int size : cookies) {
            if (fed < greed.length && size >= greed[fed]) {
                fed++;
            }
        }
        return fed;
    }
}
