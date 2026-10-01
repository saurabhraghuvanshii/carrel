#include <algorithm>
#include <string>
using namespace std;

int longestAfterChanges(string& text, int k) {
    int count[26] = {0};
    int best = 0, most = 0, start = 0;
    for (int i = 0; i < (int)text.size(); i++) {
        most = max(most, ++count[text[i] - 'A']);
        while (i - start + 1 - most > k) {
            count[text[start] - 'A']--;
            start++;
        }
        best = max(best, i - start + 1);
    }
    return best;
}
