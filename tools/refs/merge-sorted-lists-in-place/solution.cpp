#include <vector>
using namespace std;

void mergeInto(vector<int>& first, int m, vector<int>& second) {
    int i = m - 1, j = (int)second.size() - 1;
    for (int w = m + (int)second.size() - 1; j >= 0; w--) {
        if (i >= 0 && first[i] > second[j]) first[w] = first[i--];
        else first[w] = second[j--];
    }
}
