import java.util.*;

// Ported from the owner's Greedy/MaxLengthChain.java.
class Solution {
    public int longestChain(int[][] pairs) {
        Arrays.sort(pairs, Comparator.comparingDouble(o -> o[1]));
        int chainLn = 1;
        int chainEnd = pairs[0][1];
        for (int i = 1; i < pairs.length; i++) {
            if (pairs[i][0] > chainEnd) {
                chainLn++;
                chainEnd = pairs[i][1];
            }
        }
        return chainLn;
    }
}
