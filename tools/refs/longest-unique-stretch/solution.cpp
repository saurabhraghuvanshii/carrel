#include <algorithm>
#include <string>
#include <vector>
using namespace std;

int longestUnique(string& text) {
    vector<int> last(128, -1);
    int best = 0, start = 0;
    for (int i = 0; i < (int)text.size(); i++) {
        start = max(start, last[(unsigned char)text[i]] + 1);
        last[(unsigned char)text[i]] = i;
        best = max(best, i - start + 1);
    }
    return best;
}
